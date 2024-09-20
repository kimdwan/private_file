import "./assets/css/User.css"
import { UserGoMain, UserNickName, UserProfile } from "./components"

export const User = ({ computerNumber, setComputerNumber }) => {
  return (
    <div className = "userContainer">
      
      {/* 유저 프로필 */}
      <UserProfile computerNumber={ computerNumber } setComputerNumber={ setComputerNumber } />

      {/* 유저 닉네임 */}
      <UserNickName computerNumber={ computerNumber } setComputerNumber={ setComputerNumber } />

      {/* 유저 메인 페이지 */}
      <UserGoMain computerNumber={ computerNumber } setComputerNumber={ setComputerNumber } />

    </div>
  )
}