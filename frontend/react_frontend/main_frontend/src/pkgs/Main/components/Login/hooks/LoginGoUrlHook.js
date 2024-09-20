import { useCallback } from "react"
import { useNavigate } from "react-router-dom"

export const useLoginGoUrlHook = () => {
  const navigate = useNavigate()

  const clickGoSignUpBtn = useCallback((event) => {
    if (event.target.className === "loginFooterGoUrlBoxButton") {
      navigate("/signup/term")
    }

  }, [ navigate ])

  return { clickGoSignUpBtn }
}