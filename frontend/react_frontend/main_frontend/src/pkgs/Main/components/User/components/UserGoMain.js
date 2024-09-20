import { useUserGoMainUrlHook } from "../hooks"


export const UserGoMain = ({ computerNumber, setComputerNumber }) => {
  const { clickUserGoMainBtn } = useUserGoMainUrlHook(computerNumber, setComputerNumber)

  return (
    <div className = "userGoMainContainer">
      {/* 메인 화면과 로그아웃으로 이동하게 해줌 */}
      <div className = "userGoMainSmallBox">
        <button className = "userGoMainMainBtn" onClick = { clickUserGoMainBtn }>메인</button>
        <button className = "userGoMainLogoutBtn" onClick = { clickUserGoMainBtn }>로그아웃</button>
      </div>
    </div>
  )
}