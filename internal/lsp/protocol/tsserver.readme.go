package protocol

import (
	"context"
	"encoding/json"

	"golang.org/x/tools/internal/jsonrpc2"
	errors "golang.org/x/xerrors"
)

/*
	把路由部分摘出来；方便浏览

	含义参考 https://microsoft.github.io/language-server-protocol/specification
*/

func serverDispatch(ctx context.Context, server Server, reply jsonrpc2.Replier, r jsonrpc2.Request) (bool, error) {
	switch r.Method() {
	case "workspace/didChangeWorkspaceFolders": // notif
		/*
			不知道干嘛的，先跳过
		*/
		var params DidChangeWorkspaceFoldersParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidChangeWorkspaceFolders(ctx, &params)
		return true, reply(ctx, nil, err)
	case "window/workDoneProgress/cancel": // notif
		/*
			通知服务端；关闭因为 window/workDoneProgress/create 初始化的进程
		*/
		var params WorkDoneProgressCancelParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.WorkDoneProgressCancel(ctx, &params)
		return true, reply(ctx, nil, err)
	case "workspace/didCreateFiles": // notif
		/*
			通知服务端；文件创建
		*/
		var params CreateFilesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidCreateFiles(ctx, &params)
		return true, reply(ctx, nil, err)
	case "workspace/didRenameFiles": // notif
		/*
			通知服务端；文件重命名
		*/
		var params RenameFilesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidRenameFiles(ctx, &params)
		return true, reply(ctx, nil, err)
	case "workspace/didDeleteFiles": // notif
		/*
			通知服务端；文件删除
		*/
		var params DeleteFilesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidDeleteFiles(ctx, &params)
		return true, reply(ctx, nil, err)
	case "initialized": // notif
		/*
			通知服务端；客户端创建时，会发起一个 initialize 请求给服务端；然后再告知服务端已经初始化完成
		*/
		var params InitializedParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.Initialized(ctx, &params)
		return true, reply(ctx, nil, err)
	case "exit": // notif
		/*
			通知服务端；前置动作是 客户端 发起 shutdown 请求；shutdown 之后才能发起 exit 通知
			shutdown -> 不服务； exit -> 进程退出
		*/
		err := server.Exit(ctx)
		return true, reply(ctx, nil, err)
	case "workspace/didChangeConfiguration": // notif
		var params DidChangeConfigurationParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidChangeConfiguration(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/didOpen": // notif
		/*
			通知服务端；文档打开
		*/
		var params DidOpenTextDocumentParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidOpen(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/didChange": // notif
		/*
			通知服务端；文档更改
		*/
		var params DidChangeTextDocumentParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidChange(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/didClose": // notif
		/*
			通知服务端；文档关闭
		*/
		var params DidCloseTextDocumentParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidClose(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/didSave": // notif
		var params DidSaveTextDocumentParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidSave(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/willSave": // notif
		var params WillSaveTextDocumentParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.WillSave(ctx, &params)
		return true, reply(ctx, nil, err)
	case "workspace/didChangeWatchedFiles": // notif
		var params DidChangeWatchedFilesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.DidChangeWatchedFiles(ctx, &params)
		return true, reply(ctx, nil, err)
	case "$/setTrace": // notif
		var params SetTraceParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.SetTrace(ctx, &params)
		return true, reply(ctx, nil, err)
	case "$/logTrace": // notif
		var params LogTraceParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		err := server.LogTrace(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/implementation": // req
		var params ImplementationParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Implementation(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/typeDefinition": // req
		var params TypeDefinitionParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.TypeDefinition(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/documentColor": // req
		var params DocumentColorParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.DocumentColor(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/colorPresentation": // req
		var params ColorPresentationParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.ColorPresentation(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/foldingRange": // req
		var params FoldingRangeParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.FoldingRange(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/declaration": // req
		var params DeclarationParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Declaration(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/selectionRange": // req
		var params SelectionRangeParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.SelectionRange(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/prepareCallHierarchy": // req
		var params CallHierarchyPrepareParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.PrepareCallHierarchy(ctx, &params)
		return true, reply(ctx, resp, err)
	case "callHierarchy/incomingCalls": // req
		var params CallHierarchyIncomingCallsParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.IncomingCalls(ctx, &params)
		return true, reply(ctx, resp, err)
	case "callHierarchy/outgoingCalls": // req
		var params CallHierarchyOutgoingCallsParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.OutgoingCalls(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/semanticTokens/full": // req
		var params SemanticTokensParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.SemanticTokensFull(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/semanticTokens/full/delta": // req
		var params SemanticTokensDeltaParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.SemanticTokensFullDelta(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/semanticTokens/range": // req
		var params SemanticTokensRangeParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.SemanticTokensRange(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/semanticTokens/refresh": // req
		if len(r.Params()) > 0 {
			return true, reply(ctx, nil, errors.Errorf("%w: expected no params", jsonrpc2.ErrInvalidParams))
		}
		err := server.SemanticTokensRefresh(ctx)
		return true, reply(ctx, nil, err)
	case "textDocument/linkedEditingRange": // req
		var params LinkedEditingRangeParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.LinkedEditingRange(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/willCreateFiles": // req
		var params CreateFilesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.WillCreateFiles(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/willRenameFiles": // req
		var params RenameFilesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.WillRenameFiles(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/willDeleteFiles": // req
		var params DeleteFilesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.WillDeleteFiles(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/moniker": // req
		var params MonikerParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Moniker(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/prepareTypeHierarchy": // req
		var params TypeHierarchyPrepareParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.PrepareTypeHierarchy(ctx, &params)
		return true, reply(ctx, resp, err)
	case "typeHierarchy/supertypes": // req
		var params TypeHierarchySupertypesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Supertypes(ctx, &params)
		return true, reply(ctx, resp, err)
	case "typeHierarchy/subtypes": // req
		var params TypeHierarchySubtypesParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Subtypes(ctx, &params)
		return true, reply(ctx, resp, err)
	case "initialize": // req
		var params ParamInitialize
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			if _, ok := err.(*json.UnmarshalTypeError); !ok {
				return true, sendParseError(ctx, reply, err)
			}
		}
		resp, err := server.Initialize(ctx, &params)
		return true, reply(ctx, resp, err)
	case "shutdown": // req
		if len(r.Params()) > 0 {
			return true, reply(ctx, nil, errors.Errorf("%w: expected no params", jsonrpc2.ErrInvalidParams))
		}
		err := server.Shutdown(ctx)
		return true, reply(ctx, nil, err)
	case "textDocument/willSaveWaitUntil": // req
		var params WillSaveTextDocumentParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.WillSaveWaitUntil(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/completion": // req
		var params CompletionParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Completion(ctx, &params)
		return true, reply(ctx, resp, err)
	case "completionItem/resolve": // req
		var params CompletionItem
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Resolve(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/hover": // req
		var params HoverParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Hover(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/signatureHelp": // req
		var params SignatureHelpParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.SignatureHelp(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/definition": // req
		var params DefinitionParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Definition(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/references": // req
		var params ReferenceParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.References(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/documentHighlight": // req
		var params DocumentHighlightParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.DocumentHighlight(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/documentSymbol": // req
		var params DocumentSymbolParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.DocumentSymbol(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/codeAction": // req
		var params CodeActionParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.CodeAction(ctx, &params)
		return true, reply(ctx, resp, err)
	case "codeAction/resolve": // req
		var params CodeAction
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.ResolveCodeAction(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/symbol": // req
		var params WorkspaceSymbolParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Symbol(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/codeLens": // req
		var params CodeLensParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.CodeLens(ctx, &params)
		return true, reply(ctx, resp, err)
	case "codeLens/resolve": // req
		var params CodeLens
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.ResolveCodeLens(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/codeLens/refresh": // req
		if len(r.Params()) > 0 {
			return true, reply(ctx, nil, errors.Errorf("%w: expected no params", jsonrpc2.ErrInvalidParams))
		}
		err := server.CodeLensRefresh(ctx)
		return true, reply(ctx, nil, err)
	case "textDocument/documentLink": // req
		var params DocumentLinkParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.DocumentLink(ctx, &params)
		return true, reply(ctx, resp, err)
	case "documentLink/resolve": // req
		var params DocumentLink
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.ResolveDocumentLink(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/formatting": // req
		var params DocumentFormattingParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Formatting(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/rangeFormatting": // req
		var params DocumentRangeFormattingParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.RangeFormatting(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/onTypeFormatting": // req
		var params DocumentOnTypeFormattingParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.OnTypeFormatting(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/rename": // req
		var params RenameParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Rename(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/prepareRename": // req
		var params PrepareRenameParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.PrepareRename(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/executeCommand": // req
		var params ExecuteCommandParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.ExecuteCommand(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/diagnostic": // req
		var params string
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.Diagnostic(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/diagnostic": // req
		var params WorkspaceDiagnosticParams
		if err := json.Unmarshal(r.Params(), &params); err != nil {
			return true, sendParseError(ctx, reply, err)
		}
		resp, err := server.DiagnosticWorkspace(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/diagnostic/refresh": // req
		if len(r.Params()) > 0 {
			return true, reply(ctx, nil, errors.Errorf("%w: expected no params", jsonrpc2.ErrInvalidParams))
		}
		err := server.DiagnosticRefresh(ctx)
		return true, reply(ctx, nil, err)

	default:
		return false, nil
	}
}
