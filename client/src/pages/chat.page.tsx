import { useChats } from '@/contexts/ChatCtx';
import { useState } from 'react';

export const ChatPage = () => {
  const { friends } = useChats();
  const [selectedFriend, setSelectedFriend] = useState(null); // Track selected friend for chat
  const [message, setMessage] = useState(''); // State to hold the new message

  // Display a fallback message when no friends are available
  if (!friends || friends.length === 0) {
    return (
      <div className='h-full w-full max-w-screen-2xl flex justify-center items-center bg-gray-900 text-gray-300'>
        <p>No friends found. Start adding some!</p>
      </div>
    );
  }

  const handleSelectFriend = (friend) => {
    setSelectedFriend(friend); // Set selected friend when clicked
  };

  const handleSendMessage = () => {
    if (message.trim() === '') return; // Prevent sending empty messages

    // For now, just log the message to console
    console.log('Message Sent:', message, 'To:', selectedFriend.name);

    setMessage(''); // Reset message input field
  };

  return (
    <div className='h-full w-full max-w-screen-2xl flex bg-gray-900 text-gray-300 overflow-hidden'>
      {/* Friends List */}
      <div className='flex-[0.4] h-full border-r border-gray-700 p-4 overflow-y-auto'>
        <h2 className='text-xl font-bold mb-4 text-gray-100'>Friends</h2>
        {friends.map((friend) => (
          <div
            key={friend.id}
            onClick={() => handleSelectFriend(friend)}
            className='p-2 hover:bg-gray-800 rounded transition cursor-pointer'
          >
            <p className='text-sm font-medium text-gray-200'>{friend.name}</p>
            <p className='text-xs text-gray-500'>{friend.email}</p>
          </div>
        ))}
      </div>

      {/* Chat Window */}
      <div className='flex-1 h-full flex flex-col p-4'>
        {selectedFriend ? (
          <>
            {/* Chat Header */}
            <div className='flex items-center border-b border-gray-700 pb-2 mb-4'>
              <h3 className='text-xl font-bold text-gray-100'>
                Chat with {selectedFriend.name}
              </h3>
            </div>

            {/* Chat Messages */}
            <div className='flex-1 overflow-y-auto mb-4 p-2 bg-gray-800 rounded-lg h-full'>
              {/* Sample Chat Messages */}
              <div className='flex flex-col space-y-2'>
                <div className='self-start bg-gray-700 text-gray-300 p-2 rounded-lg max-w-xs'>
                  <p>Hello, how are you?</p>
                </div>
                <div className='self-end bg-blue-600 text-gray-100 p-2 rounded-lg max-w-xs'>
                  <p>I'm good, how about you?</p>
                </div>
                {/* More messages can be mapped here */}
              </div>
            </div>

            {/* Message Input */}
            <div className='flex items-center border-t border-gray-700 pt-2'>
              <input
                type='text'
                className='flex-1 bg-gray-700 text-gray-300 p-2 rounded-lg'
                placeholder='Type a message...'
                value={message}
                onChange={(e) => setMessage(e.target.value)}
              />
              <button
                className='ml-2 bg-blue-600 text-gray-100 p-2 rounded-lg'
                onClick={handleSendMessage}
              >
                Send
              </button>
            </div>
          </>
        ) : (
          <div className='flex-1 flex justify-center items-center'>
            <p className='text-gray-500'>Select a friend to start chatting.</p>
          </div>
        )}
      </div>
    </div>
  );
};

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
