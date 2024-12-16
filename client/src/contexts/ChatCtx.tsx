import { axiosInstance } from '@/lib/axiosInstance';
import {
  createContext,
  ReactElement,
  useContext,
  useEffect,
  useState,
} from 'react';

type TUser = {
  name: string;
  email: string;
  id: string;
  created_at: string;
  updated_at: string;
};

type FRIEND_SHIPS = {
  id: string;
  created_at: string;
  updated_at: string;
  friend: TUser;
  messages: TMessage[];
}[];

type INIT_STATE = FRIEND_SHIPS;

type TMessage = {
  id: string;
  text: string;
  friendship_id: string;
  from: string;
  last_id: number;
  sent: boolean;
  seen: boolean;
  created_at: string;
  updated_at: string;
};

type ChatContextType = {
  friends: TUser[];

  isSettingUp: boolean;
  getMessage: (friendship_id: string) => TMessage[] | [];
  messages: Record<string, TMessage[]>;
};

const ChatCtx = createContext<ChatContextType>({
  friends: [],

  isSettingUp: false,

  messages: {},
  getMessage: (id: string) => [],
});

export const ChatsCtxProvider = ({ children }: { children: ReactElement }) => {
  const [isSettingUp, setIsSettingUp] = useState(false);
  const [friends, setFriends] = useState<TUser[]>([]);
  const [messages, setMessages] = useState<Record<string, TMessage[]>>({});

  useEffect(() => {
    (async () => {
      try {
        setIsSettingUp(true);
        const { data } = await axiosInstance.get<INIT_STATE>('/init-state');
        const sorted_data = data.sort((b, a) => {
          const aDate = a.messages[0]?.created_at
            ? new Date(a.messages[0].created_at).getTime()
            : 0; // Use 0 (or some fallback timestamp) if messages[0] is undefined
          const bDate = b.messages[0]?.created_at
            ? new Date(b.messages[0].created_at).getTime()
            : 0;
          return aDate - bDate;
        });

        setFriends(sorted_data.map((d) => d.friend));
        let messages: Record<string, TMessage[]> = {};

        sorted_data.forEach((d) => {
          messages[d.id] = d.messages;
        });

        setMessages(messages);
      } catch (error) {
        console.error('Error initializing friends:', error);
      } finally {
        setIsSettingUp(false);
      }
    })();
  }, []);

  const addMessage = (friend_id: string, message: TMessage) => {};

  const removeMessage = (friend_id: string, message: TMessage) => {};

  const deleteFriend = (friendship_id: string) => {};

  const cleanChat = (friendship_id: string) => {};

  const getMessage = async (friendship_id: string) => {
    const message = messages[friendship_id];

    if (message) return message;

    try {
      const { data } = await axiosInstance.get('/messages');
      console.log('data: ', data);
    } catch {}
  };
  return (
    <ChatCtx.Provider value={{ friends, isSettingUp, messages, getMessage }}>
      {children}
    </ChatCtx.Provider>
  );
};

export const useChats = () => useContext(ChatCtx);
