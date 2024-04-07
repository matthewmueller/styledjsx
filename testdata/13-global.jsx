import SlackButton from "./slack-button"
import Button from "./button"

export default function Header({ success }) {
  return (
    <header>
      <nav className="buttons">
        {/* <Button href={`mailto:${props.email}`}>Contact</Button> */}
        <Button href="/faq">FAQ</Button>
        <SlackButton success={success} />
      </nav>
      <style jsx>
        {`
          .buttons {
            display: flex;
            align-items: center;
            position: absolute;
            right: 10px;
            top: 20px;
          }

          .buttons > :global(*) {
            margin-right: 15px;
          }

          .buttons > :global(*:last-child) {
            margin-right: 0px;
          }
        `}
      </style>
    </header>
  )
}
