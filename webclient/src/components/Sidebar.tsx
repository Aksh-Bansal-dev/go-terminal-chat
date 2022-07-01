import { useEffect, useState } from "react";
import { server } from "../config/server";
import { getColor } from "../utils/colors";

type onlineUser = {
  color: number;
  username: string;
};

export const Sidebar = () => {
  const [onlineUsers, setOnlineUsers] = useState<onlineUser[]>([]);
  const getInitialUsers = async () => {
    const res = await fetch("http://" + server + "/online-users");
    const data = await res.json();
    setOnlineUsers(data);
  };
  useEffect(() => {
    getInitialUsers();
  }, []);
  return (
    <div className="drawer drawer-mobile relative w-20 md:w-80">
      <input
        id="my-drawer-2"
        type="checkbox"
        className="drawer-toggle bg-red-200"
      />
      <div className="drawer-content flex flex-col items-center justify-center ">
        <label
          htmlFor="my-drawer-2"
          className="rounded cursor-pointer text-white text-5xl p-1 drawer-button md:hidden absolute top-0 left-3"
        >
          {"="}
        </label>
      </div>
      <div className="drawer-side">
        <label htmlFor="my-drawer-2" className="drawer-overlay"></label>
        <ul className="menu p-4 overflow-y-auto w-80 bg-primary">
          <li className="font-bold  text-white">Online</li>
          {onlineUsers.map((user, i) => (
            <li key={i} style={{ color: `#${getColor(user.color)}` }}>
              {user.username}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};
