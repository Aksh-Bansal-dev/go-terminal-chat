export const Chat = () => {
  const arr = Array(50).fill(0);
  return (
    <div className="w-full max-h-screen overflow-y-hidden">
      <div className="overflow-y-scroll h-[90%] text-secondary pl-8 pt-5">
        <ul>
          {arr.map((_e, i) => (
            <li key={i}>
              Lorem ipsum dolor sit amet consectetur adipisicing elit.
              Reprehenderit nesciunt ratione qui alias perspiciatis illo? Quam
              numquam fugit consequuntur tempore beatae, minus doloremque
              dolorum ut similique, odit sunt odio quo?
            </li>
          ))}
        </ul>
      </div>
      <div className="p-4">
        <input
          type="text"
          className="rounded w-full h-10 bg-primary text-secondary px-2 text-lg"
        />
      </div>
    </div>
  );
};
