/**
 * AST Parser for React component analysis
 */

import * as parser from '@typescript-eslint/parser';
import { AST_NODE_TYPES, TSESTree } from '@typescript-eslint/typescript-estree';
import { PropDefinition } from '../models/Component';

export interface ComponentInfo {
  name: string;
  type: 'functional' | 'class' | 'arrow';
  props: PropDefinition[];
  composedOf?: string[];
  isForwardRef?: boolean;
}

export interface ParseResult {
  isComponent: boolean;
  componentInfo?: ComponentInfo;
  dependencies: string[];
  jsxElements?: string[];
  elementUsageCounts?: Record<string, number>;
  hasConditionalRendering?: boolean;
  hasListRendering?: boolean;
  hooks?: string[];
  customHooks?: string[];
  usesCompoundPattern?: boolean;
  usesRenderProps?: boolean;
  isHOC?: boolean;
  error?: {
    type: string;
    message: string;
  };
}

export class ASTParser {
  private cache: Map<string, ParseResult> = new Map();

  /**
   * Parse a file and extract component information
   */
  parse(content: string, filePath: string): ParseResult {
    // Check cache
    const cacheKey = `${filePath}:${content.length}`;
    if (this.cache.has(cacheKey)) {
      return this.cache.get(cacheKey)!;
    }

    try {
      // Parse with TypeScript ESLint parser
      const ast = parser.parse(content, {
        ecmaVersion: 2022,
        sourceType: 'module',
        ecmaFeatures: {
          jsx: true
        },
        loc: true,
        range: true
      });

      const result = this.analyzeAST(ast, content, filePath);

      // Cache result
      this.cache.set(cacheKey, result);

      return result;
    } catch (error: any) {
      return {
        isComponent: false,
        dependencies: [],
        error: {
          type: error.name || 'SyntaxError',
          message: error.message
        }
      };
    }
  }

  /**
   * Analyze the AST to extract component information
   */
  private analyzeAST(ast: TSESTree.Program, content: string, filePath: string): ParseResult {
    const result: ParseResult = {
      isComponent: false,
      dependencies: [],
      jsxElements: [],
      elementUsageCounts: {},
      hooks: [],
      customHooks: []
    };

    // Track imports
    const imports = this.extractImports(ast);
    result.dependencies = imports;

    // Find component declarations
    const component = this.findComponent(ast);
    if (component) {
      result.isComponent = true;
      result.componentInfo = component;

      // Check for patterns
      result.isHOC = this.isHigherOrderComponent(ast);
      result.usesCompoundPattern = this.usesCompoundPattern(ast, content);
      result.usesRenderProps = this.usesRenderProps(ast);
    }

    // Extract JSX usage
    const jsxAnalysis = this.analyzeJSX(ast);
    result.jsxElements = jsxAnalysis.elements;
    result.elementUsageCounts = jsxAnalysis.counts;
    result.hasConditionalRendering = jsxAnalysis.hasConditional;
    result.hasListRendering = jsxAnalysis.hasListRendering;

    // Extract hooks
    const hooksAnalysis = this.analyzeHooks(ast);
    result.hooks = hooksAnalysis.reactHooks;
    result.customHooks = hooksAnalysis.customHooks;

    // Set composed components
    if (result.componentInfo && result.jsxElements) {
      result.componentInfo.composedOf = result.jsxElements
        .filter(el => el[0] === el[0].toUpperCase()); // Component names start with uppercase
    }

    return result;
  }

  /**
   * Extract imports from AST
   */
  private extractImports(ast: TSESTree.Program): string[] {
    const imports: string[] = [];

    ast.body.forEach(node => {
      if (node.type === AST_NODE_TYPES.ImportDeclaration) {
        imports.push(node.source.value as string);
      }
    });

    return imports;
  }

  /**
   * Find React component in AST
   */
  private findComponent(ast: TSESTree.Program): ComponentInfo | undefined {
    for (const node of ast.body) {
      // Check for function component
      if (node.type === AST_NODE_TYPES.FunctionDeclaration) {
        if (this.isFunctionalComponent(node)) {
          return this.extractFunctionalComponent(node);
        }
      }

      // Check for exported variable declaration (arrow function)
      if (node.type === AST_NODE_TYPES.ExportNamedDeclaration) {
        // Check for export function Component()
        if (node.declaration?.type === AST_NODE_TYPES.FunctionDeclaration) {
          if (this.isFunctionalComponent(node.declaration)) {
            return this.extractFunctionalComponent(node.declaration);
          }
        }

        if (node.declaration?.type === AST_NODE_TYPES.VariableDeclaration) {
          const varDecl = node.declaration.declarations[0];
          if (varDecl?.init?.type === AST_NODE_TYPES.ArrowFunctionExpression ||
              varDecl?.init?.type === AST_NODE_TYPES.FunctionExpression) {
            return this.extractArrowComponent(varDecl);
          }

          // Check for forwardRef
          if (this.isForwardRefComponent(varDecl)) {
            return this.extractForwardRefComponent(varDecl);
          }
        }

        // Check for class component
        if (node.declaration?.type === AST_NODE_TYPES.ClassDeclaration) {
          if (this.isClassComponent(node.declaration)) {
            return this.extractClassComponent(node.declaration);
          }
        }
      }

      // Check for default export
      if (node.type === AST_NODE_TYPES.ExportDefaultDeclaration) {
        if (node.declaration.type === AST_NODE_TYPES.FunctionDeclaration) {
          if (this.isFunctionalComponent(node.declaration)) {
            return this.extractFunctionalComponent(node.declaration);
          }
        }
        // Also check for default export of arrow functions
        if (node.declaration.type === AST_NODE_TYPES.ArrowFunctionExpression) {
          return this.extractDefaultArrowComponent(node.declaration);
        }
      }
    }

    return undefined;
  }

  /**
   * Check if node is a functional component
   */
  private isFunctionalComponent(node: TSESTree.FunctionDeclaration): boolean {
    // Has JSX return
    return this.hasJSXReturn(node.body);
  }

  /**
   * Check if node has JSX return
   */
  private hasJSXReturn(node: TSESTree.BlockStatement | TSESTree.Expression | null): boolean {
    if (!node) return false;

    if (node.type === AST_NODE_TYPES.BlockStatement) {
      return node.body.some(stmt => {
        if (stmt.type === AST_NODE_TYPES.ReturnStatement && stmt.argument) {
          return this.isJSXExpression(stmt.argument);
        }
        return false;
      });
    }

    return this.isJSXExpression(node);
  }

  /**
   * Check if expression contains JSX
   */
  private isJSXExpression(node: TSESTree.Expression): boolean {
    // Direct JSX
    if (node.type === AST_NODE_TYPES.JSXElement ||
        node.type === AST_NODE_TYPES.JSXFragment) {
      return true;
    }

    // Parenthesized JSX - Note: TypeScript AST might not have this as separate type
    // The parser usually unwraps parentheses, so JSX in parens would be direct JSX

    // Conditional JSX
    if (node.type === AST_NODE_TYPES.ConditionalExpression) {
      return this.isJSXExpression(node.consequent) ||
             this.isJSXExpression(node.alternate);
    }

    // Logical expressions (e.g., condition && <Component />)
    if (node.type === AST_NODE_TYPES.LogicalExpression) {
      return this.isJSXExpression(node.right);
    }

    // Call expressions (e.g., React.createElement)
    if (node.type === AST_NODE_TYPES.CallExpression) {
      const callee = node.callee;
      if (callee.type === AST_NODE_TYPES.MemberExpression) {
        const obj = callee.object as any;
        const prop = callee.property as any;
        if (obj.name === 'React' && prop.name === 'createElement') {
          return true;
        }
      }
    }

    return false;
  }

  /**
   * Extract functional component info
   */
  private extractFunctionalComponent(node: TSESTree.FunctionDeclaration): ComponentInfo {
    const name = node.id?.name || 'Anonymous';
    const props = this.extractPropsFromParams(node.params);

    return {
      name,
      type: 'functional',
      props
    };
  }

  /**
   * Extract arrow component info
   */
  private extractArrowComponent(node: TSESTree.VariableDeclarator): ComponentInfo {
    const name = (node.id as TSESTree.Identifier).name;
    const func = node.init as TSESTree.ArrowFunctionExpression | TSESTree.FunctionExpression;
    const props = this.extractPropsFromParams(func.params);

    return {
      name,
      type: 'arrow',
      props
    };
  }

  /**
   * Extract default arrow component info
   */
  private extractDefaultArrowComponent(node: TSESTree.ArrowFunctionExpression): ComponentInfo {
    const props = this.extractPropsFromParams(node.params);

    return {
      name: 'default',
      type: 'functional',
      props
    };
  }

  /**
   * Check if component is forwardRef
   */
  private isForwardRefComponent(node: TSESTree.VariableDeclarator): boolean {
    if (node.init?.type === AST_NODE_TYPES.CallExpression) {
      const callee = node.init.callee;
      if (callee.type === AST_NODE_TYPES.MemberExpression) {
        const obj = callee.object as TSESTree.Identifier;
        const prop = callee.property as TSESTree.Identifier;
        return obj.name === 'React' && prop.name === 'forwardRef';
      }
    }
    return false;
  }

  /**
   * Extract forwardRef component info
   */
  private extractForwardRefComponent(node: TSESTree.VariableDeclarator): ComponentInfo {
    const name = (node.id as TSESTree.Identifier).name;

    return {
      name,
      type: 'functional',
      props: [],
      isForwardRef: true
    };
  }

  /**
   * Check if class is React component
   */
  private isClassComponent(node: TSESTree.ClassDeclaration): boolean {
    if (node.superClass) {
      if (node.superClass.type === AST_NODE_TYPES.MemberExpression) {
        const obj = node.superClass.object as TSESTree.Identifier;
        const prop = node.superClass.property as TSESTree.Identifier;
        return obj.name === 'React' && prop.name === 'Component';
      }
    }
    return false;
  }

  /**
   * Extract class component info
   */
  private extractClassComponent(node: TSESTree.ClassDeclaration): ComponentInfo {
    const name = node.id?.name || 'Anonymous';

    return {
      name,
      type: 'class',
      props: []
    };
  }

  /**
   * Extract props from function parameters
   */
  private extractPropsFromParams(params: TSESTree.Parameter[]): PropDefinition[] {
    const props: PropDefinition[] = [];

    if (params.length > 0) {
      const firstParam = params[0];

      // Handle destructured props
      if (firstParam.type === AST_NODE_TYPES.ObjectPattern) {
        firstParam.properties.forEach(prop => {
          if (prop.type === AST_NODE_TYPES.Property) {
            const key = prop.key as TSESTree.Identifier;
            props.push({
              name: key.name,
              type: 'any', // Would need type analysis for accurate type
              required: true,
              description: '',
              examples: [],
              defaultValue: undefined
            });
          }
        });
      }
    }

    return props;
  }

  /**
   * Analyze JSX usage in the component
   */
  private analyzeJSX(ast: TSESTree.Program): {
    elements: string[];
    counts: Record<string, number>;
    hasConditional: boolean;
    hasListRendering: boolean;
  } {
    const elements: string[] = [];
    const counts: Record<string, number> = {};
    let hasConditional = false;
    let hasListRendering = false;

    const visit = (node: any) => {
      if (node.type === AST_NODE_TYPES.JSXElement) {
        const opening = node.openingElement;
        if (opening.name.type === AST_NODE_TYPES.JSXIdentifier) {
          const name = opening.name.name;
          if (!elements.includes(name)) {
            elements.push(name);
          }
          counts[name] = (counts[name] || 0) + 1;
        }
      }

      // Check for conditional rendering
      if (node.type === AST_NODE_TYPES.ConditionalExpression ||
          (node.type === AST_NODE_TYPES.LogicalExpression && node.operator === '&&')) {
        hasConditional = true;
      }

      // Check for list rendering (map)
      if (node.type === AST_NODE_TYPES.CallExpression) {
        if (node.callee.type === AST_NODE_TYPES.MemberExpression) {
          const prop = node.callee.property as TSESTree.Identifier;
          if (prop.name === 'map') {
            hasListRendering = true;
          }
        }
      }

      // Recurse through children
      for (const key in node) {
        if (key !== 'parent' && node[key]) {
          if (Array.isArray(node[key])) {
            node[key].forEach(visit);
          } else if (typeof node[key] === 'object') {
            visit(node[key]);
          }
        }
      }
    };

    visit(ast);

    return {
      elements,
      counts,
      hasConditional,
      hasListRendering
    };
  }

  /**
   * Analyze React hooks usage
   */
  private analyzeHooks(ast: TSESTree.Program): {
    reactHooks: string[];
    customHooks: string[];
  } {
    const reactHooks: string[] = [];
    const customHooks: string[] = [];

    const visit = (node: any) => {
      if (node.type === AST_NODE_TYPES.CallExpression) {
        let hookName: string | null = null;

        // Check for direct hook calls (useState, useEffect, etc.)
        if (node.callee.type === AST_NODE_TYPES.Identifier) {
          const name = node.callee.name;
          if (name.startsWith('use')) {
            hookName = name;
          }
        }
        // Check for React.useState, React.useEffect, etc.
        else if (node.callee.type === AST_NODE_TYPES.MemberExpression) {
          if (node.callee.object.type === AST_NODE_TYPES.Identifier &&
              node.callee.object.name === 'React' &&
              node.callee.property.type === AST_NODE_TYPES.Identifier &&
              node.callee.property.name.startsWith('use')) {
            hookName = node.callee.property.name;
          }
        }

        if (hookName) {
          const isReactHook = [
            'useState', 'useEffect', 'useContext', 'useReducer',
            'useCallback', 'useMemo', 'useRef', 'useImperativeHandle',
            'useLayoutEffect', 'useDebugValue'
          ].includes(hookName);

          if (isReactHook) {
            if (!reactHooks.includes(hookName)) {
              reactHooks.push(hookName);
            }
          } else {
            if (!customHooks.includes(hookName)) {
              customHooks.push(hookName);
            }
          }
        }
      }

      // Recurse
      for (const key in node) {
        if (key !== 'parent' && node[key]) {
          if (Array.isArray(node[key])) {
            node[key].forEach(visit);
          } else if (typeof node[key] === 'object') {
            visit(node[key]);
          }
        }
      }
    };

    visit(ast);

    return {
      reactHooks,
      customHooks
    };
  }

  /**
   * Check if component is a Higher-Order Component
   */
  private isHigherOrderComponent(ast: TSESTree.Program): boolean {
    for (const node of ast.body) {
      if (node.type === AST_NODE_TYPES.ExportNamedDeclaration) {
        const decl = node.declaration;
        if (decl?.type === AST_NODE_TYPES.VariableDeclaration) {
          const varDecl = decl.declarations[0];
          if (varDecl?.init?.type === AST_NODE_TYPES.ArrowFunctionExpression) {
            const arrow = varDecl.init;
            // HOCs typically return a function
            if (arrow.body.type === AST_NODE_TYPES.ArrowFunctionExpression ||
                arrow.body.type === AST_NODE_TYPES.FunctionExpression) {
              return true;
            }
            // Check if it's a block statement that returns a function
            if (arrow.body.type === AST_NODE_TYPES.BlockStatement) {
              for (const stmt of arrow.body.body) {
                if (stmt.type === AST_NODE_TYPES.ReturnStatement &&
                    (stmt.argument?.type === AST_NODE_TYPES.ArrowFunctionExpression ||
                     stmt.argument?.type === AST_NODE_TYPES.FunctionExpression)) {
                  return true;
                }
              }
            }
          }
        }
      }
    }
    return false;
  }

  /**
   * Check if component uses compound pattern
   */
  private usesCompoundPattern(ast: TSESTree.Program, content: string): boolean {
    // Check for Card.Header, Card.Body pattern in JSX elements
    // This pattern is common in compound components
    return /<\w+\.\w+/.test(content);
  }

  /**
   * Check if component uses render props
   */
  private usesRenderProps(ast: TSESTree.Program): boolean {
    // Look for props with names like 'render' or functions passed as children
    const visit = (node: any): boolean => {
      if (node.type === AST_NODE_TYPES.JSXAttribute) {
        const name = node.name?.name;
        if (name === 'render' || name === 'children') {
          if (node.value?.type === AST_NODE_TYPES.JSXExpressionContainer) {
            const expr = node.value.expression;
            if (expr.type === AST_NODE_TYPES.ArrowFunctionExpression ||
                expr.type === AST_NODE_TYPES.FunctionExpression) {
              return true;
            }
          }
        }
      }

      for (const key in node) {
        if (key !== 'parent' && node[key]) {
          if (Array.isArray(node[key])) {
            for (const child of node[key]) {
              if (typeof child === 'object' && visit(child)) {
                return true;
              }
            }
          } else if (typeof node[key] === 'object' && visit(node[key])) {
            return true;
          }
        }
      }

      return false;
    };

    return visit(ast);
  }
}