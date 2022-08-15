import { ChakraProvider } from "@chakra-ui/react";
import { Provider } from "jotai";
import React from "react";
import { BrowserRouter as Router } from "react-router-dom";

import Loader from "./components/loader/Loader";
import Routings from "./router/Routings";
import { customTheme } from "./styles/customTheme";

// TODO fix theme
const App = () => (
  <ChakraProvider resetCSS theme={customTheme}>
    <React.Suspense fallback={<Loader />}>
      <Provider>
        <Router>
          <Routings />
        </Router>
      </Provider>
    </React.Suspense>
  </ChakraProvider>
);

export default App;
