import profileImg from "../assets/img/profileImg.png"
import { useUserProfileGetProfileHook } from "../hooks"


export const UserProfile = ({ computerNumber, setComputerNumber }) => {
  const { userProfile } = useUserProfileGetProfileHook( computerNumber, setComputerNumber )

  return (
    <div className = "userProfileContainer">
      <img  className = "userProfileImage" src = { userProfile ? userProfile :  profileImg } alt = "프로필 이미지" />
    </div>
  )
}