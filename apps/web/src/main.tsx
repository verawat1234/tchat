
  import React from "react";
  import { createRoot } from "react-dom/client";
  import { Provider } from "react-redux";
  import App from "./App.tsx";
  import { store } from "./store";
  import "./index.css";

  // Start MSW in development (temporarily disabled for RTK verification)
  // if (import.meta.env.DEV) {
  //   const { worker } = await import('./mocks/browser');
  //   await worker.start({
  //     onUnhandledRequest: 'bypass',
  //   });
  // }

  createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
      <Provider store={store}>
        <App />
      </Provider>
    </React.StrictMode>
  );
  