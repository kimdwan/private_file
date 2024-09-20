import { SignUpForm, SignUpTerm } from "./components"
import { useSignUpGoUrlHook } from "./hooks"

export const SignUp = () => {
  const { pathType } = useSignUpGoUrlHook()

  return (
    <div className = "signUpContainer">
      {
        pathType === "term" ? <SignUpTerm /> : pathType === "form" ? <SignUpForm /> : <></>
      }
    </div>
  )
}