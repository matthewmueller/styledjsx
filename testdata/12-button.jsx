export default function Button(props) {
  return (
    <button href={props.href}>
      {props.children}
      <style jsx>
        {`
          button {
            outline: none;
            height: 33px;
            border: 1px solid #ccced1;
            text-decoration: none;
            background: white;
            border-radius: 6px;
            padding: 6px 11px 8px 11px;
            display: inline-block;
            user-select: none;
            color: #3b434b;
            font-size: 15px;
            font-weight: 700;
            cursor: pointer;
          }

          button:hover {
            border-color: #aaa;
          }

          button:active {
            border-color: #aaa;
          }
        `}
      </style>
    </button>
  )
}
