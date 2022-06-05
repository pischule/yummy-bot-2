import Item from "./Item";

function ItemList(props) {

  return (
    <div className="items-list">
      {props.items.map((item) => (
        <Item
          key={item.id}
          id={item.id}
          name={item.name}
          quantity={item.quantity}
          updateQuantity={props.updateQuantity}
        />
      ))}
    </div>
  );
}

export default ItemList;
