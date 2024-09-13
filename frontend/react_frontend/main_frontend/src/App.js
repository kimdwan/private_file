import { BrowserRouter as Routers, Routes, Route } from "react-router-dom"
import { Main } from "./pkgs";

function App() {
  return (
    <div className="App">
      <Routers>
        <Routes>
          <Route path = "/" element = {<Main />} />
        </Routes>
      </Routers>
    </div>
  );
}

export default App;
