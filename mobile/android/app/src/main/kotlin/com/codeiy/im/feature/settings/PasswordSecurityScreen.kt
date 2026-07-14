package com.codeiy.im.feature.settings

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import com.codeiy.im.core.auth.SessionStore
import com.codeiy.im.feature.login.PasswordResetScreen
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PasswordSecurityScreen(session: SessionStore, onBack: () -> Unit) {
    var showReset by remember { mutableStateOf(false) }
    var currentPassword by remember { mutableStateOf("") }
    var newPassword by remember { mutableStateOf("") }
    var confirmPassword by remember { mutableStateOf("") }
    var loading by remember { mutableStateOf(false) }
    var error by remember { mutableStateOf<String?>(null) }
    val scope = rememberCoroutineScope()

    if (showReset) {
        PasswordResetScreen(
            service = session.api.service,
            initialEmail = session.currentUser?.email.orEmpty(),
            onBack = { showReset = false },
        )
        return
    }

    val canSubmit = currentPassword.isNotEmpty() && newPassword.length >= 10 &&
        newPassword == confirmPassword && !loading

    Scaffold(topBar = { TopAppBar(title = { Text("账户安全") }, navigationIcon = { TextButton(onClick = onBack) { Text("返回") } }) }) { padding ->
        Column(Modifier.padding(padding).padding(24.dp).fillMaxSize()) {
            Text(session.currentUser?.email.orEmpty(), color = MaterialTheme.colorScheme.onSurfaceVariant)
            if (session.currentUser?.hasPassword == true) {
                OutlinedTextField(currentPassword, { currentPassword = it }, label = { Text("当前密码") }, singleLine = true, visualTransformation = PasswordVisualTransformation(), modifier = Modifier.fillMaxWidth().padding(top = 20.dp))
                OutlinedTextField(newPassword, { newPassword = it }, label = { Text("新密码（至少 10 个字符）") }, singleLine = true, visualTransformation = PasswordVisualTransformation(), modifier = Modifier.fillMaxWidth().padding(top = 12.dp))
                OutlinedTextField(confirmPassword, { confirmPassword = it }, label = { Text("确认新密码") }, singleLine = true, visualTransformation = PasswordVisualTransformation(), modifier = Modifier.fillMaxWidth().padding(top = 12.dp))
                Button(
                    onClick = {
                        loading = true
                        error = null
                        scope.launch {
                            try { session.changePassword(currentPassword, newPassword) }
                            catch (e: Exception) { error = e.message ?: "密码修改失败" }
                            finally { loading = false }
                        }
                    },
                    enabled = canSubmit,
                    modifier = Modifier.fillMaxWidth().padding(top = 20.dp),
                ) { if (loading) CircularProgressIndicator() else Text("修改密码") }
                Text("修改后所有设备需要重新登录。", style = MaterialTheme.typography.bodySmall, modifier = Modifier.padding(top = 8.dp))
            } else {
                Text("当前账户通过第三方登录。设置密码前需要验证登录邮箱。", modifier = Modifier.padding(top = 20.dp))
                Button(onClick = { showReset = true }, modifier = Modifier.padding(top = 16.dp)) { Text("验证邮箱并设置密码") }
            }
            error?.let { Text(it, color = MaterialTheme.colorScheme.error, modifier = Modifier.padding(top = 12.dp)) }
        }
    }
}
