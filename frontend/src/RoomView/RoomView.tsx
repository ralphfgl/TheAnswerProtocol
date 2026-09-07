import './RoomView.css'

function RoomView({ data, talk, onCommand }) {
    if (!data || !data.room) {
        return <p>Loading...</p>
    }
    return (
        <section className='room_view'>
            <h2 className='room_title'>Current room</h2>
            <h3>{data.room.name}</h3>
            <p>{data.room.description}</p>
            <hr />
            <p>Move:</p>
            <div className='room_buttons'>
                {Object.keys(data.room.exits).map((value, index) => (
                    <button className='room_button' key={index} onClick={() => onCommand("MOVE " + value.toUpperCase())}>{value.toUpperCase()}</button>
                ))}
            </div>
            <p>Items:</p>
            {data.items
                ?
                <ul className='items_list'>
                    {(data.items.map((value, index) => (
                        <li key={index}>
                            <div className='item_item'>
                                <p>{value}</p>
                                <button className='item_button' onClick={() => onCommand("TAKE " + value)}>TAKE</button>
                            </div>
                        </li>
                    )))}
                </ul>
                :
                <p>There is no items</p>
            }
            <hr />
            <p>NPCs: </p>
            {data.npcs.map((value, index) => (
                <div key={index}>
                    <p>{value}</p>
                    <div className='room_actions'>
                        <button className='room_button' onClick={() => onCommand("TALK " + value)}>TALK</button>
                        <button className='room_button' onClick={() => onCommand("ATTACK " + value)}>ATTACK</button>
                    </div>
                </div>
            ))}
            <hr />
            {talk && <p>{talk.dialogue}</p>}
            {/* <form className='room_form'>
                <input className='room_input' type="text" placeholder='Chat...' />
                <button className='room_button'>Send</button>
            </form> */}
        </section>
    )
}

export default RoomView
