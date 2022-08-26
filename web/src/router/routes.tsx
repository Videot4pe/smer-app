import { lazy } from "react";
import type { PathRouteProps } from "react-router-dom";

import Layout from "../layout";

const Signin = lazy(() => import("../pages/auth/Signin"));
const Signup = lazy(() => import("../pages/auth/Signup"));
const ResetPassword = lazy(() => import("../pages/auth/ResetPassword"));
const ChangePassword = lazy(() => import("../pages/auth/ChangePassword"));
const Smers = lazy(() => import("../pages/smers/Smers"));
const Profile = lazy(() => import("../pages/profile/Profile"));

export const routes: Array<PathRouteProps> = [
  {
    path: "/signin",
    element: <Signin />,
  },
  {
    path: "/signup",
    element: <Signup />,
  },
  {
    path: "/reset-password",
    element: <ResetPassword />,
  },
  {
    path: "/change-password/:hash",
    element: <ChangePassword />,
  },
];

export const privateRoutes: Array<PathRouteProps> = [
  {
    path: "/",
    element: (
      <Layout>
        <Smers />
      </Layout>
    ),
  },
  {
    path: "/smers",
    element: (
      <Layout>
        <Smers />
      </Layout>
    ),
  },
  {
    path: "/profile",
    element: (
      <Layout>
        <Profile />
      </Layout>
    ),
  },
];
