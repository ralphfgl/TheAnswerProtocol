import './Header.css'

function Header({ data, nickname, onLogout }) {
  return (
    <header className='header'>
      <p className='header_item'>Name: {nickname}</p>
      <p className='header_item'>HP: {data.hp}/{data.max_hp}</p>
      <p className='header_item'>Players on the server: {data.players || null}</p>
      <p className='header_item'>Players in the room: {data.players_room || null}</p>
      <p className='header_item'>Status: {data.status}</p>
      <button className='header_button' onClick={onLogout}>Quit</button>
    </header>
  )
}

export default Header
