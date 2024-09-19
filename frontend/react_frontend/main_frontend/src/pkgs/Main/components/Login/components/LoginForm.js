import { useLoginFormHook } from "../hooks"

export const LoginForm = ({ setComputerNumber }) => {
  const { register, handleSubmit, errors, onSubmit } = useLoginFormHook(setComputerNumber)

  return (
    <div className="loginFormContainer">
      <form onSubmit={handleSubmit(onSubmit)}>

        {/* 이메일 입력 필드 */}
        <div className="loginFormEmailDivBox">
          <div className="loginFormEmailSmallBox">
            <label className="loginFormEmailLabel" htmlFor="loginFormEmailInput">이메일: </label>
            <input
              id="loginFormEmailInput"
              {...register("email")}
              type="text" 
            />
          </div>

          <div className="loginFormEmailErrorBox">
            {errors.email?.message && <p className="loginFormErrorMsg">{errors.email.message}</p>}
          </div>
        </div>

        {/* 비밀번호 입력 필드 */}
        <div className="loginFormPasswordDivBox">
          <div className="loginFormPasswordSmallBox">
            <label className="loginFormPasswordLabel" htmlFor="loginFormPasswordInput">비밀번호: </label>
            <input
              id="loginFormPasswordInput"
              {...register("password")}
              type="password"  
            />
          </div>

          <div className="loginFormPasswordErrorBox">
            {errors.password?.message && <p className="loginFormErrorMsg">{errors.password.message}</p>}
          </div>
        </div>

        {/* 로그인 버튼 */}
        <div className="loginFormSubmitButton">
          <input
            className="loginFormSubmitInputButton"
            type="submit"
            value="로그인"
          />
        </div>
      </form>
    </div>
  )
}
