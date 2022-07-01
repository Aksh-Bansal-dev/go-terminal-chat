export const Sidebar = () => {
  return (
    <div className="drawer drawer-mobile relative w-80">
      <input
        id="my-drawer-2"
        type="checkbox"
        className="drawer-toggle bg-red-200"
      />
      <div className="drawer-content flex flex-col items-center justify-center ">
        <label
          htmlFor="my-drawer-2"
          className="rounded cursor-pointer text-white text-5xl p-1 drawer-button lg:hidden absolute top-0 left-3"
        >
          {"="}
        </label>
      </div>
      <div className="drawer-side">
        <label htmlFor="my-drawer-2" className="drawer-overlay"></label>
        <ul className="menu p-4 overflow-y-auto w-80 text-white bg-primary">
          <li className="font-bold">Online</li>
          <li>Sidebar Item 1</li>
          <li>Sidebar Item 2</li>
        </ul>
      </div>
    </div>
  );
};
