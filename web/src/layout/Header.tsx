import "./Layout.scss";
import { Box, Flex, Heading, IconButton } from "@chakra-ui/react";
import { useAtom } from "jotai";
import { ImExit, ImProfile } from "react-icons/all";
import { Link, useNavigate } from "react-router-dom";

import { jwtToken } from "../store";

import ThemeToggle from "./ThemeToggle";

const Header = () => {
  const navigate = useNavigate();
  const [, setToken] = useAtom(jwtToken);
  const handleToSignin = () => navigate("/signin");

  const logOut = () => {
    setToken(undefined);
    handleToSignin();
  };

  return (
    <Flex
      as="header"
      width="full"
      align="center"
      alignSelf="flex-start"
      justifyContent="center"
      gridGap={2}
    >
      <Link to="/">
        <Heading className="multicolor" as="h1" size="sm">
          SMER
        </Heading>
      </Link>

      <Box marginLeft="auto">
        <ThemeToggle />
        <Link to="/profile">
          <IconButton ml={2} aria-label="profile" icon={<ImProfile />} />
        </Link>
        <IconButton
          ml={2}
          aria-label="exit"
          icon={<ImExit />}
          onClick={logOut}
        />
      </Box>
    </Flex>
  );
};

export default Header;
