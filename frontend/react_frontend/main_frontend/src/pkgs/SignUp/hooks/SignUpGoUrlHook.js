import { useEffect, useState } from "react"
import { useLocation } from "react-router-dom"

export const useSignUpGoUrlHook = () => {
  const location = useLocation()
  const [ pathType, setPathType ] = useState("")
  useEffect(() => {
    const urlName = location.pathname
    const urlList = urlName.split("/")
    const pathName = urlList[urlList.length - 1]
    setPathType(pathName)
  } , [ location ])

  return { pathType }
}