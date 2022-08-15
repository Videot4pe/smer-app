import { lazy } from "react";
import type { PathRouteProps } from "react-router-dom";

import Layout from "../layout";

const Signin = lazy(() => import("../pages/auth/Signin"));
const Signup = lazy(() => import("../pages/auth/Signup"));
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
