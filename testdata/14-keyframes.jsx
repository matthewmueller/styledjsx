import SlackButton from "./slack-button"
import Button from "./button"

export default function Header({ success }) {
  return (
    <header>
      <div className="scene">
        <div className="cloud-1" />
        <div className="cloud-2" />
      </div>
      <style jsx>
        {`
          .cloud-1 {
            background-image: url("/static/images/cloud-1.svg");
            position: absolute;
            width: 159px;
            height: 92px;
            top: 50px;
            right: 800px;
            animation: float infinite 5s ease-in;
          }

          .cloud-2 {
            background-image: url("/static/images/cloud-2.svg");
            position: absolute;
            height: 128px;
            width: 228px;
            right: 500px;
            top: 30px;
            animation: float infinite 5s ease-in;
          }

          @keyframes float {
            0% {
              transform: translate(0, 0);
            }

            50% {
              transform: translate(7px, 0px) rotate(5deg);
            }
          }
        `}
      </style>
    </header>
  )
}
