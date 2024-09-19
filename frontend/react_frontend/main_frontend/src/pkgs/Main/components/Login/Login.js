import "./assets/css/Login.css"

import { LoginFooter, LoginForm, LoginLogo } from "./components"

export const Login = ({ setComputerNumber }) => {
  return (
    <div className = "loginContainer">
      
      {/* 로그인의 로고가 들어가는 장소 */}
      <LoginLogo />

      {/* 로그인이 이루어지는 내용 */}
      <LoginForm setComputerNumber = { setComputerNumber } />

      {/* 로그인에 발 */}
      <LoginFooter />

    </div>
  )
}