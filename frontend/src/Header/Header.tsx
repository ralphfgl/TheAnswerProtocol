import './Header.css'

function Header({data, onLogout}) {
  return (
    <header className='header'>
      <p className='header_item'>Name: {data.name}</p>
      <p className='header_item'>HP: {data.health}/100</p>
      <p className='header_item'>Players on the server: {data.players}</p>
      <p className='header_item'>Players in the room: {data.players_room}</p>
      <button className='header_button' onClick={onLogout}>Quit</button>
    </header>
  )
}

export default Header
