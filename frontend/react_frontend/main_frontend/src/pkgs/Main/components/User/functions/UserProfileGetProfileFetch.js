

export const UserProfileGetProfileFetch = async ( url, computerNumber, setComputerNumber ) => {

  try {
    const response = await fetch(url, {
      method : "GET",
      headers : {
        "Content-Type" : "application/json; charset=utf-8",
        "User-Computer-Number" : computerNumber,
        "X-Requested-With" : "XMLHttpRequest",
      },
      credentials : "include",
    })

    if (!response.ok) {
      if (response.status === 401) {
        localStorage.removeItem("logan_computer_number")
        setComputerNumber("")
        alert("세션이 만료 되었습니다")
        window.location.href = "/"
        throw new Error("세션 만료")
      } else if (response.status === 500) {
        alert("서버에 오류가 발생했습니다")
        throw new Error("서버에 오류가 발생했습니다")
      } else {
        alert("오류가 발생했습니다")
        throw new Error(`오류가 발생했습니다 오류번호: ${response.status}`)
      }
    }

    const data = await response.json()
    if (data && data["imagetype"] && data["imagebase64"]) {
      const imageData = `data:image/${data["imagetype"]};base64,${data["imagebase64"]}`
      return imageData
    } else {
      return undefined
    }
  } catch (err) {
    throw err
  }

}