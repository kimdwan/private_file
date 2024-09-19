import "./assets/css/User.css"
import { UserGoMain, UserNickName, UserProfile } from "./components"

export const User = ({ computerNumber }) => {
  return (
    <div className = "userContainer">
      
      {/* 유저 프로필 */}
      <UserProfile />

      {/* 유저 닉네임 */}
      <UserNickName />

      {/* 유저 메인 페이지 */}
      <UserGoMain />

    </div>
  )
}