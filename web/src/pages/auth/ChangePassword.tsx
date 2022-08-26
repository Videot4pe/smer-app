import {
    Box,
    Heading,
    Text,
  } from "@chakra-ui/react";
  
  import Card from "./components/Card";
  import Logo from "./components/Logo";
  import ChangePasswordForm from "./components/ChangePasswordForm";
import { useParams } from "react-router-dom";
  
  const ChangePassword = () => {
    const { hash } = useParams();

    return (<Box bg="pink.100" minH="100vh" py="12" px={{ base: "4", lg: "8" }}>
      <Box maxW="md" mx="auto">
        <Logo mx="auto" h="8" mb={{ base: "10", md: "20" }} />
        <Heading textAlign="center" size="xl" fontWeight="extrabold">
          Reset password
        </Heading>
        <Text mt="4" mb="8" align="center" maxW="md" fontWeight="medium">
          <Text as="span">Set your new password</Text>
        </Text>
        <Card>
          <ChangePasswordForm hash={ hash! }/>
        </Card>
      </Box>
    </Box>
  )};
  
  export default ChangePassword;
  