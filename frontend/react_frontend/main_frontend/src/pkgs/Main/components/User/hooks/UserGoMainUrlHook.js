import { useCallback } from "react"
import { useNavigate } from "react-router-dom"
import { UserGoMainLogoutFetch } from "../functions"

export const useUserGoMainUrlHook = ( computerNumber, setComputerNumber ) => {
  const navigate = useNavigate()

  const clickUserGoMainBtn = useCallback( async (event) => {
    const eventClassName = event.target.className

    if (computerNumber) {
      if (eventClassName === "userGoMainMainBtn") {
        navigate("/main/list/")
      } else if (eventClassName === "userGoMainLogoutBtn") {
        const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL
        const url = `${go_backend_url}/auth/logout`
        const response = await UserGoMainLogoutFetch(url, computerNumber, setComputerNumber)
        if (response) {
          localStorage.removeItem("logan_computer_number")
          setComputerNumber("")
          alert(response["message"])
          navigate("/")
        }
      }
    }

  }, [ computerNumber, setComputerNumber, navigate ])

  return { clickUserGoMainBtn }
}