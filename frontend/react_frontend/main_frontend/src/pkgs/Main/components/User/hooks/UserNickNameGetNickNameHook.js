import { useState, useEffect } from "react"
import { UserNickNameGetNickNameFetch } from "../functions"

export const useUserNickNameHook = ( computerNumber, setComputerNumber ) => {
  const [ userNickName, setUserNickName ] = useState("")
  useEffect(() => {
    const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL
    const url = `${go_backend_url}/auth/getnickname`
    const getUserNickNameFunc = async (url, computerNumber, setComputerNumber) => {
      const response = await UserNickNameGetNickNameFetch(url, computerNumber, setComputerNumber)
      if (response) {
        setUserNickName(response["nickname"])
      }
    }

    if (computerNumber) {
      getUserNickNameFunc(url, computerNumber, setComputerNumber)
    }

  }, [ computerNumber, setComputerNumber ])

  return { userNickName }
}