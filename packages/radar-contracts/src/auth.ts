import type { components } from './openapi';

export type AuthToken = components['schemas']['AuthTokenRead'];
type GeneratedAuthUser = components['schemas']['AuthUserRead'];
export interface AuthUser extends GeneratedAuthUser {
  hasPassword: boolean;
}
export type NativeLoginRequest = components['schemas']['NativeLoginRequest'];
export type PasswordLoginRequest = components['schemas']['PasswordLoginRequest'];
