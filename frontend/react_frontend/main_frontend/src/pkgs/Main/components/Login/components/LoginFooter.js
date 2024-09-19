import { useLoginGoUrlHook } from "../hooks"


export const LoginFooter = () => {
  const { clickGoSignUpBtn } = useLoginGoUrlHook()

  return (
    <div className = "loginFooterContainer">
      <div className = "loginFooterGoUrlBox">
        <button onClick = {clickGoSignUpBtn} className = "loginFooterGoUrlBoxButton">
          회원가입 하러 가기
        </button>
      </div>
    </div>
  )
}