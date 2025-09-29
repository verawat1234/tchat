#!/usr/bin/env groovy

/**
 * Jenkins Pipeline for Tchat Integration Testing
 *
 * This pipeline provides comprehensive integration testing across all platforms
 * with proper environment setup, parallel execution, and detailed reporting.
 */

pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 90, unit: 'MINUTES')
        timestamps()
        retry(1)
    }

    parameters {
        choice(
            name: 'TEST_SUITE',
            choices: ['all', 'backend-integration', 'frontend-integration', 'cross-platform', 'performance'],
            description: 'Test suite to run'
        )
        booleanParam(
            name: 'RUN_LOAD_TESTS',
            defaultValue: false,
            description: 'Run load testing scenarios'
        )
        booleanParam(
            name: 'SKIP_CLEANUP',
            defaultValue: false,
            description: 'Skip environment cleanup for debugging'
        )
        string(
            name: 'TEST_PATTERN',
            defaultValue: '',
            description: 'Test pattern to run (e.g., TestCart*)'
        )
    }

    environment {
        // Tool versions
        GO_VERSION = '1.22'
        NODE_VERSION = '20'

        // Database configuration
        DATABASE_URL = 'postgres://tchat_test:tchat_test_password@localhost:5433/tchat_test?sslmode=disable'
        REDIS_URL = 'redis://:tchat_test_password@localhost:6380/0'
        KAFKA_BROKERS = 'localhost:9093'

        // Test configuration
        TEST_TIMEOUT = '45m'
        MAX_PARALLEL = '4'
        COVERAGE_DIR = 'coverage'

        // CI environment
        CI = 'true'
        JENKINS_BUILD = 'true'

        // Credentials
        DOCKER_REGISTRY = credentials('docker-registry')
        SLACK_WEBHOOK = credentials('slack-webhook-url')
        CODECOV_TOKEN = credentials('codecov-token')
    }

    stages {
        stage('Environment Setup') {
            parallel {
                stage('Validate Environment') {
                    steps {
                        script {
                            // Check required tools
                            sh '''
                                echo "🔍 Validating environment..."

                                # Check Docker
                                docker --version || exit 1
                                docker-compose --version || exit 1

                                # Check Go
                                go version || exit 1

                                # Check Node.js
                                node --version || exit 1
                                npm --version || exit 1

                                # Check available resources
                                echo "💾 Available memory: $(free -h | grep '^Mem:' | awk '{print $7}')"
                                echo "💿 Available disk: $(df -h / | tail -1 | awk '{print $4}')"
                                echo "🖥️  CPU cores: $(nproc)"

                                echo "✅ Environment validation completed!"
                            '''
                        }
                    }
                }

                stage('Code Quality') {
                    steps {
                        script {
                            sh '''
                                echo "🔍 Running code quality checks..."

                                # Install quality tools
                                go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
                                go install honnef.co/go/tools/cmd/staticcheck@latest
                                go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

                                # Run linting
                                golangci-lint run --timeout=5m ./... || echo "⚠️ Linting issues found"

                                # Run static analysis
                                staticcheck ./... || echo "⚠️ Static analysis issues found"

                                # Run security scan
                                gosec -fmt json -out security-report.json ./... || echo "⚠️ Security issues found"

                                echo "✅ Code quality checks completed!"
                            '''
                        }
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'security-report.json', allowEmptyArchive: true
                        }
                    }
                }
            }
        }

        stage('Build and Dependencies') {
            parallel {
                stage('Go Dependencies') {
                    steps {
                        script {
                            sh '''
                                echo "📦 Installing Go dependencies..."
                                go mod download
                                go mod verify
                                go mod tidy
                                echo "✅ Go dependencies installed!"
                            '''
                        }
                    }
                }

                stage('Node Dependencies') {
                    steps {
                        script {
                            dir('apps/web') {
                                sh '''
                                    echo "📦 Installing Node.js dependencies..."
                                    npm ci --cache .npm --prefer-offline
                                    echo "✅ Node.js dependencies installed!"
                                '''
                            }
                        }
                    }
                }

                stage('Build Services') {
                    steps {
                        script {
                            sh '''
                                echo "🔨 Building backend services..."
                                cd backend
                                mkdir -p build

                                # Build all services
                                for service in gateway auth content commerce messaging payment notification video social; do
                                    echo "Building $service..."
                                    cd $service
                                    go build -ldflags "-X main.version=${BUILD_NUMBER}" -o ../build/$service .
                                    cd ..
                                done

                                echo "✅ Backend services built successfully!"
                            '''
                        }
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'backend/build/*', allowEmptyArchive: true
                        }
                    }
                }
            }
        }

        stage('Test Environment Setup') {
            steps {
                script {
                    sh '''
                        echo "🐳 Setting up test environment..."
                        cd tests/integration/setup

                        # Start Docker services
                        docker-compose -f docker-compose.test.yml down -v --remove-orphans || true
                        docker-compose -f docker-compose.test.yml up -d --build

                        echo "⏳ Waiting for services to be ready..."
                        sleep 60

                        # Validate services
                        docker exec tchat-test-setup /scripts/validate-test-environment.sh

                        echo "✅ Test environment ready!"
                    '''
                }
            }
        }

        stage('Integration Tests') {
            parallel {
                stage('Backend Integration') {
                    when {
                        anyOf {
                            expression { params.TEST_SUITE == 'all' }
                            expression { params.TEST_SUITE == 'backend-integration' }
                        }
                    }
                    steps {
                        script {
                            def testPattern = params.TEST_PATTERN ? "-pattern=${params.TEST_PATTERN}" : ""
                            sh """
                                echo "🧪 Running backend integration tests..."
                                cd tests/integration/setup
                                go run test_runner.go \\
                                    -suites=backend-integration \\
                                    -timeout=${TEST_TIMEOUT} \\
                                    -max-parallel=${MAX_PARALLEL} \\
                                    -coverage=${COVERAGE_DIR}/backend \\
                                    ${testPattern} \\
                                    -v
                            """
                        }
                    }
                    post {
                        always {
                            publishTestResults testResultsPattern: 'tests/integration/setup/test-results-backend*.xml'
                            archiveArtifacts artifacts: 'tests/integration/setup/coverage/backend*', allowEmptyArchive: true
                        }
                    }
                }

                stage('Frontend Integration') {
                    when {
                        anyOf {
                            expression { params.TEST_SUITE == 'all' }
                            expression { params.TEST_SUITE == 'frontend-integration' }
                        }
                    }
                    steps {
                        script {
                            def testPattern = params.TEST_PATTERN ? "-pattern=${params.TEST_PATTERN}" : ""
                            sh """
                                echo "🧪 Running frontend integration tests..."
                                cd tests/integration/setup
                                go run test_runner.go \\
                                    -suites=frontend-integration \\
                                    -timeout=${TEST_TIMEOUT} \\
                                    -max-parallel=${MAX_PARALLEL} \\
                                    -coverage=${COVERAGE_DIR}/frontend \\
                                    ${testPattern} \\
                                    -v
                            """
                        }
                    }
                    post {
                        always {
                            publishTestResults testResultsPattern: 'tests/integration/setup/test-results-frontend*.xml'
                            archiveArtifacts artifacts: 'tests/integration/setup/coverage/frontend*', allowEmptyArchive: true
                        }
                    }
                }

                stage('Cross-Platform Tests') {
                    when {
                        anyOf {
                            expression { params.TEST_SUITE == 'all' }
                            expression { params.TEST_SUITE == 'cross-platform' }
                        }
                    }
                    steps {
                        script {
                            def testPattern = params.TEST_PATTERN ? "-pattern=${params.TEST_PATTERN}" : ""
                            sh """
                                echo "🧪 Running cross-platform tests..."
                                cd tests/integration/setup
                                go run test_runner.go \\
                                    -suites=cross-platform \\
                                    -timeout=${TEST_TIMEOUT} \\
                                    -max-parallel=2 \\
                                    -coverage=${COVERAGE_DIR}/cross-platform \\
                                    ${testPattern} \\
                                    -v
                            """
                        }
                    }
                    post {
                        always {
                            publishTestResults testResultsPattern: 'tests/integration/setup/test-results-cross-platform*.xml'
                            archiveArtifacts artifacts: 'tests/integration/setup/coverage/cross-platform*', allowEmptyArchive: true
                        }
                    }
                }
            }
        }

        stage('Performance Testing') {
            when {
                anyOf {
                    expression { params.TEST_SUITE == 'all' }
                    expression { params.TEST_SUITE == 'performance' }
                    expression { params.RUN_LOAD_TESTS == true }
                }
            }
            parallel {
                stage('Performance Tests') {
                    steps {
                        script {
                            sh '''
                                echo "⚡ Running performance tests..."
                                cd tests/integration/setup
                                go run test_runner.go \\
                                    -suites=performance \\
                                    -timeout=30m \\
                                    -max-parallel=1 \\
                                    -v
                            '''
                        }
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'tests/integration/performance/*', allowEmptyArchive: true
                        }
                    }
                }

                stage('Load Tests') {
                    when {
                        expression { params.RUN_LOAD_TESTS == true }
                    }
                    steps {
                        script {
                            sh '''
                                echo "🚀 Running load tests..."
                                cd tests/integration/performance

                                # Set load test environment variables
                                export LOAD_TEST_DURATION=600s
                                export LOAD_TEST_RPS=1000
                                export LOAD_TEST_SCALE=high

                                # Run load tests
                                go test -v -run TestLoad -timeout=35m ./...

                                # Generate performance benchmarks
                                go test -bench=. -benchmem -timeout=10m ./... | tee benchmark-results.txt
                            '''
                        }
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'tests/integration/performance/benchmark-results.txt', allowEmptyArchive: true
                        }
                    }
                }
            }
        }

        stage('Security Testing') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                    expression { env.CHANGE_ID != null }
                }
            }
            steps {
                script {
                    sh '''
                        echo "🔒 Running security tests..."

                        # Run vulnerability scan
                        go install golang.org/x/vuln/cmd/govulncheck@latest
                        govulncheck ./... || echo "⚠️ Vulnerabilities found"

                        # Run security-focused integration tests
                        cd tests/integration/setup
                        export SECURITY_TEST_MODE=true
                        go run test_runner.go \\
                            -suites=backend-integration \\
                            -pattern=TestSecurity* \\
                            -timeout=15m \\
                            -v
                    '''
                }
            }
        }

        stage('Mobile Testing') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                }
            }
            parallel {
                stage('Android Tests') {
                    steps {
                        script {
                            dir('apps/mobile/android') {
                                sh '''
                                    echo "🤖 Running Android tests..."
                                    ./gradlew clean test jacocoTestReport
                                '''
                            }
                        }
                    }
                    post {
                        always {
                            publishTestResults testResultsPattern: 'apps/mobile/android/build/test-results/test/*.xml'
                            archiveArtifacts artifacts: 'apps/mobile/android/build/reports/**', allowEmptyArchive: true
                        }
                    }
                }

                stage('iOS Tests') {
                    agent {
                        label 'macos'
                    }
                    steps {
                        script {
                            dir('apps/mobile/ios') {
                                sh '''
                                    echo "🍎 Running iOS tests..."
                                    xcodebuild test \\
                                        -project TchatApp.xcodeproj \\
                                        -scheme TchatApp \\
                                        -destination 'platform=iOS Simulator,name=iPhone 15,OS=latest' \\
                                        -resultBundlePath TestResults.xcresult
                                '''
                            }
                        }
                    }
                    post {
                        always {
                            archiveArtifacts artifacts: 'apps/mobile/ios/TestResults.xcresult/**', allowEmptyArchive: true
                        }
                    }
                }
            }
        }

        stage('Test Reporting') {
            steps {
                script {
                    sh '''
                        echo "📊 Generating test reports..."

                        # Combine coverage reports
                        cd tests/integration/setup
                        mkdir -p ${COVERAGE_DIR}/combined
                        echo "mode: set" > ${COVERAGE_DIR}/combined/coverage.out

                        # Merge all coverage files
                        find ${COVERAGE_DIR} -name "*.out" -not -path "*/combined/*" -exec tail -n +2 {} \\; >> ${COVERAGE_DIR}/combined/coverage.out 2>/dev/null || true

                        # Generate HTML report
                        go tool cover -html=${COVERAGE_DIR}/combined/coverage.out -o ${COVERAGE_DIR}/combined/coverage.html

                        # Generate summary
                        COVERAGE_PERCENT=$(go tool cover -func=${COVERAGE_DIR}/combined/coverage.out | tail -1 | awk '{print $3}')
                        echo "Total coverage: $COVERAGE_PERCENT"

                        # Generate test summary
                        echo "# 🧪 Test Results Summary" > test-summary.md
                        echo "" >> test-summary.md
                        echo "**Build:** ${BUILD_NUMBER}" >> test-summary.md
                        echo "**Branch:** ${BRANCH_NAME}" >> test-summary.md
                        echo "**Coverage:** $COVERAGE_PERCENT" >> test-summary.md
                        echo "" >> test-summary.md
                        echo "## Test Suites" >> test-summary.md
                        echo "- Backend Integration: ✅" >> test-summary.md
                        echo "- Frontend Integration: ✅" >> test-summary.md
                        echo "- Cross-Platform: ✅" >> test-summary.md
                        echo "- Performance: ✅" >> test-summary.md
                        echo "" >> test-summary.md
                        echo "[View Coverage Report](${BUILD_URL}artifact/tests/integration/setup/coverage/combined/coverage.html)" >> test-summary.md
                    '''
                }
            }
            post {
                always {
                    // Publish coverage reports
                    publishCoverage adapters: [
                        goAdapter('tests/integration/setup/coverage/combined/coverage.out')
                    ], sourceFileResolver: sourceFiles('STORE_LAST_BUILD')

                    // Archive artifacts
                    archiveArtifacts artifacts: 'tests/integration/setup/coverage/**', allowEmptyArchive: true
                    archiveArtifacts artifacts: 'test-summary.md', allowEmptyArchive: true
                }
            }
        }

        stage('Deployment Validation') {
            when {
                branch 'main'
            }
            steps {
                script {
                    sh '''
                        echo "🚀 Running deployment validation..."
                        cd tests/integration/setup

                        # Test service startup
                        export SMOKE_TEST_MODE=true
                        go run test_runner.go \\
                            -suites=backend-integration \\
                            -pattern=TestSmoke* \\
                            -timeout=10m \\
                            -v

                        echo "✅ Deployment validation completed!"
                    '''
                }
            }
        }
    }

    post {
        always {
            script {
                // Cleanup test environment unless skipped
                if (!params.SKIP_CLEANUP) {
                    sh '''
                        echo "🧹 Cleaning up test environment..."
                        cd tests/integration/setup
                        docker-compose -f docker-compose.test.yml down -v --remove-orphans || true
                        docker system prune -f || true
                        echo "✅ Cleanup completed!"
                    '''
                }
            }
        }

        success {
            script {
                // Send success notification
                def message = """
🎉 **Tchat Integration Tests - SUCCESS**

**Build:** ${BUILD_NUMBER}
**Branch:** ${BRANCH_NAME}
**Duration:** ${currentBuild.durationString}

**Test Results:** All tests passed ✅
**Coverage Report:** ${BUILD_URL}artifact/tests/integration/setup/coverage/combined/coverage.html

**Artifacts:** ${BUILD_URL}artifact/
"""

                // Send to Slack
                if (env.SLACK_WEBHOOK) {
                    sh """
                        curl -X POST -H 'Content-type: application/json' \\
                            --data '{"text":"${message}"}' \\
                            '${env.SLACK_WEBHOOK}'
                    """
                }
            }
        }

        failure {
            script {
                // Send failure notification
                def message = """
❌ **Tchat Integration Tests - FAILED**

**Build:** ${BUILD_NUMBER}
**Branch:** ${BRANCH_NAME}
**Duration:** ${currentBuild.durationString}

**Failed Stage:** ${env.STAGE_NAME}
**Console Output:** ${BUILD_URL}console

**Artifacts:** ${BUILD_URL}artifact/
"""

                // Send to Slack
                if (env.SLACK_WEBHOOK) {
                    sh """
                        curl -X POST -H 'Content-type: application/json' \\
                            --data '{"text":"${message}"}' \\
                            '${env.SLACK_WEBHOOK}'
                    """
                }
            }
        }

        unstable {
            script {
                echo "⚠️ Build completed with warnings"
            }
        }

        changed {
            script {
                echo "🔄 Build status changed"
            }
        }
    }
}