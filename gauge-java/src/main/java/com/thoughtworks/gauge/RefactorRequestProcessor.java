package com.thoughtworks.gauge;


import main.Messages;
import org.walkmod.javalang.ASTManager;
import org.walkmod.javalang.JavaParser;
import org.walkmod.javalang.ast.CompilationUnit;
import org.walkmod.javalang.ast.body.*;
import org.walkmod.javalang.ast.expr.AnnotationExpr;
import org.walkmod.javalang.ast.expr.SingleMemberAnnotationExpr;
import org.walkmod.javalang.ast.expr.StringLiteralExpr;
import org.walkmod.javalang.visitors.VoidVisitorAdapter;

import java.io.File;
import java.io.FileInputStream;
import java.util.ArrayList;
import java.util.List;

public class RefactorRequestProcessor implements IMessageProcessor {
    @Override
    public Messages.Message process(Messages.Message message) {
        Messages.RefactorRequest refactorRequest = message.getRefactorRequest();
        findJavaFile(refactorRequest.getOldStepText());
        return message;
    }

    private String findJavaFile(String oldStepText) {
        File workingDir = new File(System.getProperty("user.dir"));
        List<JavaParseWorker> javaFiles = parseAllJavaFiles(workingDir);
        StepValueExtractor stepValueExtractor = new StepValueExtractor();
        for (JavaParseWorker javaFile : javaFiles) {
            CompilationUnit compilationUnit = javaFile.getCompilationUnit();
            for (TypeDeclaration typeDeclaration : compilationUnit.getTypes()) {
                if (!(typeDeclaration instanceof ClassOrInterfaceDeclaration))
                    continue;

                for (BodyDeclaration bodyDeclaration : typeDeclaration.getMembers()) {
                    if (!(bodyDeclaration instanceof MethodDeclaration))
                        continue;

                    MethodDeclaration methodDeclaration = (MethodDeclaration) bodyDeclaration;
                    for (AnnotationExpr annotationExpr : methodDeclaration.getAnnotations()) {
                        if (!(annotationExpr instanceof SingleMemberAnnotationExpr))
                            continue;

                        SingleMemberAnnotationExpr annotation = (SingleMemberAnnotationExpr) annotationExpr;
                        if (annotation.getMemberValue() instanceof StringLiteralExpr) {
                            StringLiteralExpr memberValue = (StringLiteralExpr) annotation.getMemberValue();
                            StepValueExtractor.StepValue stepValue = stepValueExtractor.getValue(memberValue.getValue());
                            System.out.println(stepValue.getValue() + "   :  " + memberValue.getValue());
                            if (stepValue.getValue().equals(oldStepText)) {
                                System.out.println("done");
                            }
                        }

                    }
                }
            }
        }

        return null;
    }

    private List<JavaParseWorker> parseAllJavaFiles(File workingDir) {
        ArrayList<JavaParseWorker> javaFiles = new ArrayList<JavaParseWorker>();
        File[] allFiles = workingDir.listFiles();
        for (File file : allFiles) {
            if (file.isDirectory()) {
                javaFiles.addAll(parseAllJavaFiles(file));
            } else {
                if (file.getName().toLowerCase().endsWith(".java")) {
                    JavaParseWorker worker = new JavaParseWorker(file);
                    worker.start();
                    javaFiles.add(worker);
                }
            }
        }

        return javaFiles;
    }

    class JavaParseWorker extends Thread {

        private File javaFile;
        private CompilationUnit compilationUnit;

        JavaParseWorker(File javaFile) {
            this.javaFile = javaFile;
        }

        @Override
        public void run() {
            try {
                FileInputStream in = new FileInputStream(javaFile);
                compilationUnit = JavaParser.parse(in);
                in.close();
            } catch (Exception e) {
                // ignore exceptions
            }
        }

        public File getJavaFile() {
            return javaFile;
        }

        CompilationUnit getCompilationUnit() {
            try {
                join();
            } catch (InterruptedException e) {

            }
            return compilationUnit;
        }
    }

}
