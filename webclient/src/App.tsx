import { Chat } from "./components/Chat";
import { Sidebar } from "./components/Sidebar";

export const App = () => {
  return (
    <div className="h-screen bg-primaryLight flex">
      <Sidebar />
      <Chat />
    </div>
  );
};
