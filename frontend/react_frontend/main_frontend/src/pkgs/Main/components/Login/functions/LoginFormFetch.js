

export const LoginFormFetch = async ( url, datas, setError ) => {

  try {
    const response = await fetch(url, {
      method : "POST",
      headers : {
        "Content-Type" : "application/json",
        "X-Requested-With" : "",
      },
      body : JSON.stringify(datas),
      credentials : "include",
    })

    if (!response.ok) {
      if (response.status === 500) {
        alert("서버에 오류가 발생했습니다")
        throw new Error("서버에 오류가 발생했습니다")
      } else if (response.status === 400) {
        alert("클라이언트 폼에 문제가 발생했습니다.")
        throw new Error("클라이언트 폼에 문제가 발생했습니다")
      } else if (response.status === 406) {
        setError("email", {
          type : "manual" ,
          message : "이메일이 존재하지 않습니다."
        })
        throw new Error("이메일이 존재하지 않습니다")
      } else if (response.status === 424) {
        setError("password", {
          type : "manual",
          message : "비밀번호가 다릅니다 다시 확인해주세요",
        })
        throw new Error("비밀번호가 다름")
      } else {
        alert("서버에 오류가 발생했습니다")
        throw new Error(`오류가 발생했습니다 에러번호: ${response.status}`)
      }
    }

    const data = await response.json()
    return data

  } catch (err) {
    throw err
  }



}