import { useUserNickNameHook } from "../hooks"

export const UserNickName = ({ computerNumber, setComputerNumber }) => {
  const { userNickName } = useUserNickNameHook(computerNumber, setComputerNumber)

  return (
    <div className = "userNickNameContainer">
      {
        userNickName
      }
    </div>
  )
}