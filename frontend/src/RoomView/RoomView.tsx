import './RoomView.css'

function RoomView({ data, onCommand }) {
    return (
        <section className='room_view'>
            <h2 className='room_title'>Current room</h2>
            <h3>{data.name}</h3>
            <p>{data.description}</p>
            <hr />
            <p>Move:</p>
            {Object.keys(data.exits).map((value, index) => (
                <button className='room_button' key={index} onClick={() => onCommand("MOVE " + value.toUpperCase())}>{value.toUpperCase()}</button>
            ))}
            <p>Items:</p>
            <p>Sword</p>
            <hr />
            <p>NPCs: </p>
            {data.spawns.map((value, index) => (
                <div key={index}>
                    <p>{value.npc_type}</p>
                    <div className='room_actions'>
                        <button className='room_button' onClick={() => onCommand("TALK " + value.npc_type)}>Talk</button>
                        <button className='room_button' onClick={() => onCommand("ATTACK " + value.npc_type)}>Attack</button>
                    </div>
                </div>
            ))}
            <hr />
            <form className='room_form'>
                <input className='room_input' type="text" placeholder='Chat...' />
                <button className='room_button'>Send</button>
            </form>
        </section>
    )
}

export default RoomView
