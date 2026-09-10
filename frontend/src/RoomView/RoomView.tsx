import './RoomView.css'

export interface RoomData {
    id: string;
    name: string;
    description: string;
    exits: Record<string, string>;
    spawns: [];
}

export interface Data {
    type: string;
    room: RoomData | null;
    npcs: string[];
    items: string[];
}

export interface Talk {
    type: string;
    npc: string;
    dialogue: string;
}

interface RoomViewProps {
    data: Data | null
    talk: Talk | null
    attack: string,
    onCommand: (command: string) => void
}

function RoomView(props: RoomViewProps) {
    const { data, talk, attack, onCommand } = props
    if (!data || !data.room) {
        return <p>Loading...</p>
    }
    return (
        <section className='room_view'>
            <h2 className='room_title'>Current room</h2>
            <h3>{data.room.name}</h3>
            <p>{data.room.description}</p>
            <hr className='line' />
            <p>Move:</p>
            <div className='room_buttons'>
                {Object.keys(data.room.exits).map((value, index) => (
                    <button className='room_button' key={index} onClick={() => onCommand("MOVE " + value)}>{value.toUpperCase()}</button>
                ))}
            </div>
            <hr className='line' />
            <p>Items:</p>
            {data.items
                ?
                (data.items.map((value, index) => (
                    <div className='item_card' key={index}>
                        <p>{value.toUpperCase()}</p>
                        <button className='item_button' onClick={() => onCommand("TAKE " + value)}>TAKE</button>
                    </div>
                )))
                :
                <p>There is no items</p>
            }
            <hr className='line' />
            <p>NPCs: </p>
            {
                data.npcs
                    ?
                    data.npcs.map((value, index) => (
                        <div className='item_card' key={index}>
                            <p>{value.toUpperCase()}</p>
                            <div className='room_actions'>
                                <button className='room_button' onClick={() => onCommand("TALK " + value)}>TALK</button>
                                <button className='room_button' onClick={() => onCommand("ATTACK " + value)}>ATTACK</button>
                                <button className='room_button' onClick={() => onCommand("FLEE")}>FLEE</button>
                                <button className='room_button' onClick={() => onCommand("QUEST " + value)}>QUEST</button>
                            </div>
                        </div>
                    ))
                    :
                    <p>No NPCs</p>
            }
            <hr className='line' />
            {talk && <p>{talk.dialogue}</p>}
            {attack && <p>{attack}</p>}
            {/* <form className='room_form'>
                <input className='room_input' type="text" placeholder='Chat...' />
                <button className='room_button'>Send</button>
            </form> */}
        </section>
    )
}

export default RoomView
