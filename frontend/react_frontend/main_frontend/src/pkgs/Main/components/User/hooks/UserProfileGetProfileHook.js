import { UserProfileGetProfileFetch } from "../functions"

import { useState, useEffect } from "react"

export const useUserProfileGetProfileHook = (computerNumber, setComputerNumber) => {
  const [ userProfile, setUserProfile ] = useState(undefined)
  useEffect(() => {
    const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL
    const url = `${go_backend_url}/auth/getprofile`
    const getProfileImageFunc = async (url, computerNumber, setComputerNumber) => {
      const response = await UserProfileGetProfileFetch(url, computerNumber, setComputerNumber)
      if (response) {
        setUserProfile(response)
      }
    }

    if (computerNumber) {
      getProfileImageFunc(url, computerNumber, setComputerNumber)
    }

  }, [ computerNumber, setComputerNumber ])

  return { userProfile }  
}