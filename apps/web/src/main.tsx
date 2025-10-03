
import React from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "react-redux";
import { RouterProvider } from "react-router-dom";
import { router } from "./router/routes";
import { store } from "./store";
import { DialogProvider } from "./components/DialogSystem";
import "./index.css";

// Import service routing tests for development
// DISABLED: Using Railway services, not localhost services
// Uncomment only when running all services locally for testing
// if (import.meta.env.DEV) {
//   import('./utils/testServiceRouting');
// }

// Start MSW in development (temporarily disabled for RTK verification)
// if (import.meta.env.DEV) {
//   const { worker } = await import('./mocks/browser');
//   await worker.start({
//     onUnhandledRequest: 'bypass',
//   });
// }

// Force light mode by removing dark class
document.documentElement.classList.remove('dark');

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider store={store}>
      <DialogProvider>
        <RouterProvider router={router} />
      </DialogProvider>
    </Provider>
  </React.StrictMode>
);
  