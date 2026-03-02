import { ResourcesConfig } from "aws-amplify";

const userPoolId = process.env.REACT_APP_COGNITO_USER_POOL_ID;
const userPoolClientId =
    process.env.REACT_APP_COGNITO_APP_CLIENT_ID ||
    process.env.REACT_APP_COGNITO_USER_POOL_CLIENT_ID;

if (!userPoolId || !userPoolClientId) {
    throw new Error(
        "Cognito configuration is missing. Set REACT_APP_COGNITO_USER_POOL_ID and REACT_APP_COGNITO_APP_CLIENT_ID (or REACT_APP_COGNITO_USER_POOL_CLIENT_ID)."
    );
}

const amplifyConfig: ResourcesConfig = {
    Auth: {
        Cognito: {
            userPoolId,
            userPoolClientId,
            loginWith: {
                email: true,
            },
        },
    },
};

export default amplifyConfig;