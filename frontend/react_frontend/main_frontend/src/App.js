import { BrowserRouter as Routers, Routes, Route } from "react-router-dom"
import { Main } from "./pkgs";
import { createContext } from "react";
import { LoadComputerNumber } from "./settings";

export const MainContext = createContext()

function App() {
  const { computerNumber, setComputerNumber } = LoadComputerNumber()

  return (
    <div className="App">
      <MainContext.Provider value = {{ computerNumber, setComputerNumber }}>
        <Routers>
          <Routes>
            <Route path = "/" element = {<Main />} />
          </Routes>
        </Routers>
      </MainContext.Provider>
    </div>
  );
}

export default App;
