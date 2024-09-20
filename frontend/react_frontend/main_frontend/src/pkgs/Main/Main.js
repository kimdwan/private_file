// 컴퍼넌트
import { User, Login } from "./components"
import { MainContext } from "../../App"

// 기본 함수
import { useContext } from "react"

export const Main = () => {
  const { computerNumber, setComputerNumber } = useContext(MainContext)

  return (
    <div>
      {
        computerNumber ? <User computerNumber = { computerNumber } setComputerNumber = { setComputerNumber } /> : <Login setComputerNumber = { setComputerNumber } />
      }
    </div>
  )
}