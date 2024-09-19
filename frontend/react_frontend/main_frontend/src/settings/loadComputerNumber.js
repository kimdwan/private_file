import { useEffect, useState } from "react"

export const LoadComputerNumber = () => {
  const [ computerNumber, setComputerNumber ] = useState(undefined)

  useEffect(() => {
    const originComputerNumber = localStorage.getItem("logan_computer_number")
    if (originComputerNumber) {
      setComputerNumber(originComputerNumber)
    }

  }, [  ])

  return { computerNumber, setComputerNumber }
}