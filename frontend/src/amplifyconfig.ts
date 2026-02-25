import { ResourcesConfig } from "aws-amplify";

const amplifyConfig: ResourcesConfig = {
    Auth: {
        Cognito: {
            userPoolId: process.env.REACT_APP_COGNITO_USER_POOL_ID || "",
            userPoolClientId: process.env.REACT_APP_COGNITO_USER_POOL_CLIENT_ID || "",
            loginWith: {
                email: true,
            },
        },
    },
};

export default amplifyConfig;