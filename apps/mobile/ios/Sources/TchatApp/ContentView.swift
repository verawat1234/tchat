import SwiftUI

struct ContentView: View {
    var body: some View {
        VStack(spacing: 12) {
            Text("Tchat iOS")
                .font(.largeTitle)
                .bold()
            Text("SwiftUI shell ready for auth and chat flows.")
                .font(.subheadline)
                .foregroundStyle(.secondary)
        }
        .padding()
    }
}

#Preview {
    ContentView()
}
