import './Header.css'

export interface HeaderData {
  type: string;
  hp: string;
  max_hp: string;
  status: string;
}

interface HeaderProps {
  data: HeaderData | null;
  nickname: string;
  onLogout: () => void;
  playersServer: number;
  playersRoom: number;
}

function Header(props: HeaderProps) {
  const { data, nickname, onLogout, playersServer, playersRoom } = props
  return (
    <header className='header'>
      <p className='header_item'>Name: {nickname}</p>
      <p className='header_item'>HP: {data?.hp}/{data?.max_hp}</p>
      <p className='header_item'>Players on the server: {playersServer}</p>
      <p className='header_item'>Players in the room: {playersRoom}</p>
      <p className='header_item'>Status: {data?.status}</p>
      <button className='header_button' onClick={onLogout}>Quit</button>
    </header>
  )
}

export default Header
