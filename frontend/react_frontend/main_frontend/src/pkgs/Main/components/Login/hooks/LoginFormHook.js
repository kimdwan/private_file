// 함수들
import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { LoginFormFetch } from "../functions"

export const useLoginFormHook = ( setComputerNumber ) => {
  const schema = yup.object({
    email : yup.string().email("이메일의 형식을 지켜주시길 바랍니다.").required("이메일은 필수적으로 입력해주셔야 합니다."),
    password : yup.string().min(4, "비밀번호는 최소4글자 입니다.").max(16, "비밀번호는 최대16글자 입니다.").required("비밀번호는 필수로 입력해야 하는 사항입니다.")
  })

  const { register, handleSubmit, formState : { errors }, setError } = useForm({
    resolver : yupResolver(schema)
  })

  const onSubmit = async (data) => {
    const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL
    const url = `${go_backend_url}/user/login`

    const response  = await LoginFormFetch(url, data, setError)
    if (response) {
      const computer_number = response["computer_number"]
      localStorage.setItem("logan_computer_number", computer_number)
      setComputerNumber(computer_number)
      alert("로그인 되었습니다.")
    }

  }
  
  return { register, handleSubmit, errors, onSubmit }
}