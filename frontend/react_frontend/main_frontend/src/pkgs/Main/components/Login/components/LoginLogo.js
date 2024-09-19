import mainLogo from "../assets/img/mainLogo.webp"

import { Link } from "react-router-dom"

export const LoginLogo = () => {
  return (
    <div className = "loginLogoContainer">

      <Link to = "/">
        <img src = { mainLogo } className = "loginLogoImage" alt = "메인로고" />
      </Link>

    </div>
  )
}