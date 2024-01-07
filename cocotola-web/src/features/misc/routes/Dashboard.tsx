import { ContentLayout } from '@/components/layout';
import { useAuthStore } from '@/stores/auth';

export const Dashboard = () => {
  const userInfo = useAuthStore((state) => state.userInfo);
  return (
    <ContentLayout title="Dashboard">
      <h1 className="text-xl mt-2">
        Welcome <b>{`${userInfo?.username}`}</b>
      </h1>
      <p className="font-medium">In this application you can:</p>
      {/* {user?.role === ROLES.USER && (
        <ul className="my-4 list-inside list-disc">
          <li>Create comments in discussions</li>
          <li>Delete own comments</li>
        </ul>
      )}
      {user?.role === ROLES.ADMIN && (
        <ul className="my-4 list-inside list-disc">
          <li>Create discussions</li>
          <li>Edit discussions</li>
          <li>Delete discussions</li>
          <li>Comment on discussions</li>
          <li>Delete all comments</li>
        </ul>
      )} */}
    </ContentLayout>
  );
};
