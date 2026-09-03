import './RoomView.css'

function RoomView({ data, onCommand }) {
    return (
        <section className='room_view'>
            <h2>Game World</h2>
            <h3>{data.name}</h3>
            <p>{data.description}</p>
            <hr />
            <p>Move:</p>
            {Object.keys(data.exits).map((value, index) => (
                <button key={index} onClick={() => onCommand("MOVE " + value.toUpperCase())}>{value}</button>
            ))}
            <p>Items:</p>
            <p>Sword</p>
            <hr />
            <p>NPCs: </p>
            {data.spawns.map((value) => (
                <div>
                    <p>{value.npc_type}</p>
                    <button>Talk</button>
                    <button>Attack</button>
                </div>
            ))}
            <hr />
            <form>
                <input type="text" placeholder='Chat...' />
                <button>Send</button>
            </form>
        </section>
    )
}

export default RoomView
