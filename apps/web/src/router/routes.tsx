import { createBrowserRouter, Navigate } from 'react-router-dom';
import { Layout } from '../components/Layout';
import { ChatTab } from '../components/ChatTab';
import { StoreTab } from '../components/StoreTab';
import { SocialTab } from '../components/SocialTab';
import { VideoTab } from '../components/VideoTab';
import { WorkspaceTab } from '../components/WorkspaceTab';
import { NotificationsScreen } from '../components/NotificationsScreen';
import { SettingsScreen } from '../components/SettingsScreen';
import { WalletScreen } from '../components/WalletScreen';
import { CartScreen } from '../components/CartScreen';
import { useGetCurrentUserQuery } from '../services/microservicesApi';

// User Provider Component that fetches real user data
const UserProvider = ({ children }: { children: (user: any) => React.ReactNode }) => {
  const { data: user, isLoading, error } = useGetCurrentUserQuery();

  // Fallback user for offline/error states
  const fallbackUser = {
    id: '1',
    name: 'Demo User',
    email: 'demo@tchat.app',
    avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=32&h=32&fit=crop&crop=face'
  };

  // Use real user data if available, fallback if not
  const currentUser = user || fallbackUser;

  return <>{children(currentUser)}</>;
};

// Define route structure
export const router = createBrowserRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      {
        index: true,
        element: <Navigate to="/chat" replace />
      },
      {
        path: 'chat',
        element: (
          <UserProvider>
            {(user) => (
              <ChatTab
                user={user}
                onSearch={() => {}}
                onNotifications={() => {}}
                onNewChat={() => {}}
                onOpenChat={() => {}}
                onAcceptCall={() => {}}
                onRejectCall={() => {}}
                onEndCall={() => {}}
                chats={[]}
                incomingCall={null}
                currentCall={null}
              />
            )}
          </UserProvider>
        )
      },
      {
        path: 'store',
        element: (
          <UserProvider>
            {(user) => (
              <StoreTab
                user={user}
                onBack={() => {}}
                onOpenProduct={() => {}}
                onAddToCart={() => {}}
                onOpenCart={() => {}}
                onOpenShop={() => {}}
                onOpenLiveStream={() => {}}
                onOpenShopChat={() => {}}
                cartItems={[]}
              />
            )}
          </UserProvider>
        )
      },
      {
        path: 'social',
        element: (
          <UserProvider>
            {(user) => (
              <SocialTab
                user={user}
                onBack={() => {}}
                onOpenPost={() => {}}
                onLikePost={() => {}}
                onCommentPost={() => {}}
                onSharePost={() => {}}
                onOpenStory={() => {}}
                onOpenProfile={() => {}}
                likedPosts={[]}
                viewedStories={[]}
              />
            )}
          </UserProvider>
        )
      },
      {
        path: 'video/*',
        element: (
          <UserProvider>
            {(user) => (
              <VideoTab
                user={user}
                onBack={() => {}}
                onVideoPlay={() => {}}
                onVideoLike={() => {}}
                onVideoShare={() => {}}
                onSubscribe={() => {}}
                currentVideoId=""
                watchedVideos={[]}
                likedVideos={[]}
                subscribedChannels={[]}
              />
            )}
          </UserProvider>
        )
      },
      {
        path: 'more',
        element: (
          <UserProvider>
            {(user) => (
              <WorkspaceTab
                user={user}
                onBack={() => {}}
                onOpenFeature={() => {}}
              />
            )}
          </UserProvider>
        )
      },
      {
        path: 'notifications',
        element: (
          <UserProvider>
            {(user) => (
              <NotificationsScreen
                user={user}
                onBack={() => {}}
                onMarkAsRead={() => {}}
                onMarkAllAsRead={() => {}}
                notifications={[]}
              />
            )}
          </UserProvider>
        )
      },
      {
        path: 'settings',
        element: (
          <UserProvider>
            {(user) => (
              <SettingsScreen
                user={user}
                onBack={() => {}}
                onUpdateProfile={() => {}}
                onUpdateSettings={() => {}}
              />
            )}
          </UserProvider>
        )
      },
      {
        path: 'wallet',
        element: (
          <UserProvider>
            {(user) => (
              <WalletScreen
                user={user}
                onBack={() => {}}
                onSendMoney={() => {}}
                onTopUp={() => {}}
                balance={0}
                transactions={[]}
              />
            )}
          </UserProvider>
        )
      },
      {
        path: 'cart',
        element: (
          <UserProvider>
            {(user) => (
              <CartScreen
                user={user}
                onBack={() => {}}
                onUpdateQuantity={() => {}}
                onRemoveItem={() => {}}
                onCheckout={() => {}}
                cartItems={[]}
              />
            )}
          </UserProvider>
        )
      }
    ]
  }
]);

export default router;